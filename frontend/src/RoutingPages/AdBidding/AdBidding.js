import H_NavBar from "../../component/navbar/H_NavBar";
import styles from "./AdBidding.module.css";
import practice from "../../component/images/preview_practice.png";
import rank from "../../component/images/preview_rank.png";
import symbol from "../../component/images/logo.png";
import previous from "../../component/images/previous.png";
import next from "../../component/images/next.png";
import { useEffect, useRef, useState } from "react";
import Countdown from "../../component/Countdown/Countdown";
import HorizontalLine from "../../component/lines/HorizontalLine";
import { useParams } from "react-router-dom";
import { BsXLg } from "react-icons/bs";
import axiosFormClient from "../../component/backendConn/axiosFormClient";
import axiosClient from "../../component/backendConn/axiosClient";
import checkAccessTokenValidity from "../../component/backendConn/checkAccessTokenValidity";

function AddBidding() {
  const { locationParam } = useParams();
  const fileInputRef = useRef(null);

  const [idx, setIdx] = useState(0);
  const titles = ["연습모드 하단", "랭크 페이지 중간", "무료 토큰 지급 페이지"];
  const previImages = [practice, rank];
  const locations = ["practice", "rank"];
  const reqImgSize = ["1500 x 70", "1060 x 70"];
  const [userID, setUserID] = useState("");
  const [highestBidAmt, setHighestBidAmt] = useState(0);
  const [nextUnlock, setNextUnlock] = useState();
  const [bidOpen, setBidOpen] = useState(false);
  const [imgPreview, setImgPreview] = useState(null);
  const [imageFileError, setImageFileError] = useState("");
  const [selectedFile, setSelectedFile] = useState(null);
  const allowedExtensions = ["jpg", "jpeg", "png", "gif"];
  const [userBidAmt, setUserBidAmt] = useState();
  const [bidAmtError, setBidAmtError] = useState("");

  const submitAd = async (e) => {
    e.preventDefault();

    const formData = new FormData();
    formData.append("location", locations[idx]);
    formData.append("amount", userBidAmt);
    formData.append("image", selectedFile);

    try {
      const response = await axiosFormClient.post("/bidToken", formData);
      console.log(response);
      if (response.status === 200) {
        setBidOpen(false);
      }
    } catch (error) {
      if (error.response.data.includes("insufficient")) {
        alert(`토큰이 부족합니다.`);
      }
      console.error(error);
    }
  };

  const highestBidder = async (location) => {
    try {
      const res = await axiosClient.get(`/highestBidder?location=${location}`);
      setUserID(res.data.user_id);
      setHighestBidAmt(res.data.amount);
    } catch {
      setUserID("아직 입찰자가 없습니다.");
      setHighestBidAmt(0);
    }
  };

  const clickPrevious = () => {
    if (idx <= 0) {
      return;
    }
    setIdx((current) => current - 1);
    const loc = locations[idx - 1];
    highestBidder(loc);
  };

  const clickNext = () => {
    if (idx >= titles.length) {
      return;
    }
    setIdx((current) => current + 1);
    const loc = locations[idx + 1];
    highestBidder(loc);
  };

  const handleFileChange = (event) => {
    const selected = event.target.files[0];
    const fileExtension = selected.name.split(".").pop().toLowerCase();
    if (!allowedExtensions.includes(fileExtension)) {
      setImageFileError(
        "이미지 파일 확장자가 잘못되었습니다. JPG, JPEG, PNG, GIF 중에서 업로드 해주세요."
      );
      return;
    }

    const maxSize = 10 * 1024 * 1024;
    if (selected.size > maxSize) {
      setImageFileError("이미지 파일은 10 MB 이내로 업로드 해주세요.");
      return;
    }

    setSelectedFile(selected);
    setImgPreview(URL.createObjectURL(selected));
    setImageFileError("");
  };

  const handleButtonClick = () => {
    fileInputRef.current.click();
  };

  const isNumber = (inputValue) => {
    return /^-?\d*\.?\d+$/.test(inputValue);
  };

  const handleAmountChange = (event) => {
    if (!isNumber(event.target.value)) {
      setBidAmtError("입찰가는 숫자만 입력 가능합니다.");
    } else {
      setBidAmtError("");
    }
    setUserBidAmt(event.target.value);
  };

  useEffect(() => {
    const getNextBidUnlock = async () => {
      const res = await axiosClient.get("/nextBidUnlock");
      setNextUnlock(res.data.next_unlock);
    };

    getNextBidUnlock();
    locations.map((loc, i) => {
      if (locationParam && loc === locationParam) {
        setIdx(i);
        highestBidder(loc);
        return;
      }
    });
    highestBidder(locations[0]);
  }, []);

  return (
    <div className={styles.adbidding}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      <div className={styles.title}>
        {(idx + 1).toString().padStart(2, "0")}
        {". "}
        {titles[idx]}
      </div>
      <div className={styles.preview}>
        <img
          className={styles.navbutton}
          src={previous}
          onClick={clickPrevious}
        />
        <img className={styles.previewimage} src={previImages[idx]} />
        <div className={styles.highestbidder}>
          <h2>최고 입찰자</h2>
          <HorizontalLine />
          <h3>{userID}</h3>
          <div className={styles.tokenbalance}>
            <img src={symbol} />
            <h3>{highestBidAmt.toLocaleString()}</h3>
          </div>
        </div>
        <img className={styles.navbutton} src={next} onClick={clickNext} />
      </div>
      <HorizontalLine />
      <div className={styles.timer}>
        <h2>입찰 마감까지</h2>
        {nextUnlock ? <Countdown nextUnlock={nextUnlock} /> : null}
        <button
          className={styles.bidbutton}
          onClick={async () => {
            const userInfo = await checkAccessTokenValidity();
            if (!userInfo) {
              alert("로그인이 필요합니다.");
              window.location.replace("/login");
              return;
            }
            setBidOpen(true);
          }}
        >
          입찰하기
        </button>
      </div>
      {bidOpen ? (
        <div className={styles.comment}>
          <div
            className={styles.background}
            onClick={() => {
              setBidOpen(false);
            }}
          ></div>
          <div className={styles.wrapper}>
            <div className={styles.closebutton}>
              <span>
                <BsXLg
                  onClick={() => {
                    setBidOpen(false);
                  }}
                />
              </span>
            </div>
            <div className={styles.title}>광고할 이미지를 등록해 주세요</div>
            <div className={styles.caution}>
              정상적인 이미지 출력을 위해 규격 ({reqImgSize[idx]})에 맞는
              이미지를 업로드 해주세요.
            </div>
            <div className={styles.imginput}>
              <img
                className={styles.imgpreview}
                src={imgPreview}
                onClick={handleButtonClick}
              />
              <button
                className={styles.selectimg}
                type="button"
                onClick={handleButtonClick}
              >
                찾아보기
              </button>
              <input
                id="image"
                type="file"
                onChange={handleFileChange}
                ref={fileInputRef}
                accept="image/*"
              />
            </div>
            <div className={styles.numberinputwrapper}>
              <input
                className={styles.numberinput}
                placeholder="광고 스팟에 대한 입찰가를 입력해주세요."
                value={userBidAmt}
                onChange={handleAmountChange}
              />
              <div className={styles.moi}>MOI</div>
            </div>
            {imageFileError || bidAmtError ? (
              <div className={styles.errormessage}>
                {imageFileError
                  ? imageFileError
                  : bidAmtError
                  ? bidAmtError
                  : null}
              </div>
            ) : null}
            <button
              className={styles.sendbutton}
              disabled={imageFileError !== "" || bidAmtError !== ""}
              onClick={submitAd}
            >
              등록하기
            </button>
          </div>
        </div>
      ) : null}
    </div>
  );
}

export default AddBidding;
