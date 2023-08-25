import axiosClient from "../../component/backendConn/axiosClient";
import H_NavBar from "../../component/navbar/H_NavBar";
import styles from "./AdBidding.module.css";
import practice from "../../component/images/preview_practice.png";
import rank from "../../component/images/preview_rank.png";
import symbol from "../../component/images/logo.png";
import previous from "../../component/images/previous.png";
import next from "../../component/images/next.png";
import { useEffect, useRef, useState } from "react";
import Countdown from "./Countdown/Countdown";
import HorizontalLine from "../../component/lines/HorizontalLine";
import { useParams } from "react-router-dom";
import { BsXLg } from "react-icons/bs";

function AddBidding() {
  const { location } = useParams();
  const fileInputRef = useRef(null);

  const [idx, setIdx] = useState(0);
  const titles = ["연습모드 하단", "랭크 페이지 중간", "무료 토큰 지급 페이지"];
  const previImages = [practice, rank];
  const pathes = ["practice", "rank"];
  const [userID, setUserID] = useState("");
  const [bidAmt, setBidAmt] = useState(0);
  const [nextUnlock, setNextUnlock] = useState();
  const [bidOpen, setBidOpen] = useState(false);
  const [imgPreview, setImgPreview] = useState(null);
  const [imageFileError, setImageFileError] = useState("");
  const [selectedFile, setSelectedFile] = useState(null);
  const allowedExtensions = ["jpg", "jpeg", "png", "gif"];

  const highestBidder = async (path) => {
    try {
      const res = await axiosClient.get(`/highestBidder?location=${path}`);
      setUserID(res.data.user_id);
      setBidAmt(res.data.amount);
    } catch (error) {
      console.error(error);
      setUserID("아직 입찰자가 없습니다.");
      setBidAmt(0);
    }
  };

  const clickPrevious = () => {
    if (idx <= 0) {
      return;
    }
    setIdx((current) => current - 1);
    const path = pathes[idx - 1];
    highestBidder(path);
  };

  const clickNext = () => {
    if (idx >= titles.length) {
      return;
    }
    setIdx((current) => current + 1);
    const path = pathes[idx + 1];
    highestBidder(path);
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

  useEffect(() => {
    const getNextBidUnlock = async () => {
      const res = await axiosClient.get("/nextBidUnlock");
      setNextUnlock(res.data.next_unlock);
    };

    getNextBidUnlock();
    pathes.map((path, i) => {
      if (location && path === location) {
        setIdx(i);
        highestBidder(path);
        return;
      }
    });
    highestBidder(pathes[0]);
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
            <h3>{bidAmt.toLocaleString()}</h3>
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
          onClick={() => {
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
          <div className={styles.inner}>
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
            <div className={styles.imginput}>
              <img className={styles.imgpreview} src={imgPreview} />
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
            <input
              className={styles.numberinput}
              type="number"
              placeholder="광고 스팟에 대한 입찰가를 입력해주세요."
              maxLength={100}
            />
            <button className={styles.sendbutton}>등록하기</button>
          </div>
        </div>
      ) : null}
    </div>
  );
}

export default AddBidding;
