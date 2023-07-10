import { useEffect, useState } from "react";
import styles from "./SignUp.module.css";
import H_NavBar from "../../component/navbar/H_NavBar";
import { LoadAccessToken } from "../../component/Token/Token";
function SignUp() {
  const userIDPattern = /^[a-zA-Z0-9]{5,15}$/;
  const passwordPattern =
    /^(?=.*[a-zA-Z])(?=.*[0-9])(?=.*[$@$!%*#?&])[a-zA-Z0-9$@$!%*#?&]{8,16}$/;
  const nicknamePattern = /^[가-힣a-zA-Z0-9]{1,10}$/;
  const [userID, setUserID] = useState("");
  const [emailID, setEmailID] = useState("");
  const [emailDomain, setEmailDomain] = useState("");
  const [selectDomainDisable, setSelectDomainDisable] = useState(true);
  const [password, setPassword] = useState("");
  const [passwordChk, setPasswordChk] = useState("");
  const [nickname, setNickname] = useState("");
  const [isLogined, setIsLogined] = useState(false);
  const [userIdDuplicationText, setUserIdDuplicationText] =
    useState("중복확인");
  const [nicknameDuplicationText, setNicknameDuplicationText] =
    useState("중복확인");
  const [userIDError, setuserIDError] = useState("");
  const [passwordError, setPasswordError] = useState("");
  const [passwordCheckError, setPasswordCheckError] = useState("");
  const [nicknameError, setNicknameError] = useState("");

  const submit = (e) => {
    e.preventDefault();
    fetch("http://bitmoi.co.kr:5000/user", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        user_id: userID,
        password: password,
        nickname: nickname,
        email: emailID + "@" + emailDomain,
        photo_url: "",
        oauth_uid: "",
      }),
    });
  };

  const userIDChange = (e) => {
    setUserID(e.target.value);
    setUserIdDuplicationText("중복확인");
    if (!userIDPattern.test(e.target.value)) {
      setuserIDError("ID는 영문 숫자 조합 5자에서 15자 내로 입력해 주세요.");
    } else {
      setuserIDError("");
    }
  };

  const userIDCheck = (e) => {
    e.preventDefault();
    fetch("http://bitmoi.co.kr:5000/user/checkid?user_id=" + userID)
      .then((res) => {
        if (res.ok) {
          setUserIdDuplicationText("사용가능");
          return;
        } else {
          setUserIdDuplicationText("사용불가");
        }
      })
      .catch((error) => {
        console.error(error);
      });
  };

  const pwChange = (e) => {
    setPassword(e.target.value);
    if (e.target.value.length < 8) {
      setPasswordError("비밀번호는 최소 8자리여야 합니다.");
    } else {
      if (!passwordPattern.test(e.target.value)) {
        setPasswordError(
          "비밀번호는 영문, 숫자, 특수문자($, @, $, !, %, *, #, ?, &)가 모두 포함되어야 합니다."
        );
      } else {
        setPasswordError("");
        if (passwordChk !== "" && e.target.value !== passwordChk) {
          setPasswordCheckError("재확인 비밀번호가 다릅니다.");
        } else {
          setPasswordCheckError("");
        }
      }
    }
  };

  const pwChkChange = (e) => {
    setPasswordChk(e.target.value);
    if (!passwordPattern.test(e.target.value)) {
      setPasswordCheckError(
        "비밀번호는 영문, 숫자, 특수문자($, @, $, !, %, *, #, ?, &)가 모두 포함되어야 합니다."
      );
    } else {
      setPasswordCheckError("");
      if (password !== e.target.value) {
        setPasswordCheckError("재확인 비밀번호가 다릅니다.");
      }
    }
  };

  const nicknameChange = (e) => {
    setNickname(e.target.value);
    setNicknameDuplicationText("중복확인");
    if (!nicknamePattern.test(e.target.value)) {
      setNicknameError(
        "닉네임은 특수문자를 제외하여 10자 이내로 입력해 주세요."
      );
    } else {
      setNicknameError("");
    }
  };

  const nicknameCheck = (e) => {
    e.preventDefault();
    fetch("http://bitmoi.co.kr:5000/user/checknickname?nickname=" + nickname)
      .then((res) => {
        if (res.ok) {
          setNicknameDuplicationText("사용가능");
          return;
        } else {
          setNicknameDuplicationText("사용불가");
        }
      })
      .catch((error) => {
        console.error(error);
      });
  };

  const emailIDChange = (e) => {
    setEmailID(e.target.value);
  };

  const selectDomain = (e) => {
    if (e.target.value === "직접입력") {
      setEmailDomain("");
      setSelectDomainDisable(false);
    } else {
      setEmailDomain(e.target.value);
      setSelectDomainDisable(true);
    }
  };
  const typingDomain = (e) => {
    setEmailDomain(e.target.value);
  };

  useEffect(() => {
    if (LoadAccessToken() === "undefined") {
      setIsLogined(false);
    } else {
      setIsLogined(true);
    }
  }, []);
  return (
    <div className={styles.signupdiv}>
      {isLogined ? (
        <div className={styles.warning}>
          <h1>잘못된 접근입니다!</h1>
          <a href="/">BACK</a>
        </div>
      ) : (
        <div className={styles.formdiv}>
          <div className={styles.navbar}>
            <H_NavBar />
          </div>
          <h3 className={styles.welcome}>
            시뮬레이션 모의투자 비트모이에 오신 걸 환영합니다!
          </h3>
          <form className={styles.forms}>
            <div className={styles.field}>
              <label htmlFor="id">아이디</label>
              <input
                id="id"
                type="text"
                placeholder="ID"
                value={userID}
                onChange={userIDChange}
                style={{
                  width: "60%",
                  backgroundColor: userIDError === "" ? "" : "#f3ebeb",
                }}
              ></input>
              <button
                className={styles.duplication}
                onClick={userIDCheck}
                style={
                  userIdDuplicationText === "사용가능"
                    ? { backgroundColor: "#3e4da3" }
                    : userIdDuplicationText === "사용불가"
                    ? { backgroundColor: "#972523bd" }
                    : { backgroundColor: "#5e5f66" }
                }
              >
                {userIdDuplicationText}
              </button>
            </div>

            <div className={styles.field}>
              <label htmlFor="pw">비밀번호</label>
              <input
                id="pw"
                type="password"
                placeholder="password"
                value={password}
                onChange={pwChange}
                style={{
                  backgroundColor: passwordError === "" ? "" : "#f3ebeb",
                }}
              ></input>
            </div>

            <div className={styles.field}>
              <label htmlFor="pwcheck">비밀번호 확인</label>
              <input
                id="pwcheck"
                type="password"
                placeholder="repeat password"
                value={passwordChk}
                onChange={pwChkChange}
                style={{
                  backgroundColor: passwordCheckError === "" ? "" : "#f3ebeb",
                }}
              ></input>
            </div>

            <div className={styles.field}>
              <label htmlFor="nickname">닉네임</label>
              <input
                id="nickname"
                type="text"
                value={nickname}
                onChange={nicknameChange}
                style={{ width: "60%" }}
              ></input>
              <button
                className={styles.duplication}
                onClick={nicknameCheck}
                style={
                  nicknameDuplicationText === "사용가능"
                    ? { backgroundColor: "#3e4da3" }
                    : nicknameDuplicationText === "사용불가"
                    ? { backgroundColor: "#972523bd" }
                    : { backgroundColor: "#5e5f66" }
                }
              >
                {nicknameDuplicationText}
              </button>
            </div>

            <div className={styles.field}>
              <label htmlFor="emailID">이메일</label>
              <input
                style={{ width: "25%" }}
                id="emailID"
                value={emailID}
                onChange={emailIDChange}
              ></input>
              <input
                style={{ width: "30%" }}
                disabled={selectDomainDisable}
                value={emailDomain}
                onChange={typingDomain}
              ></input>
              <select
                className={styles.selectBox}
                id="selectEmailDomain"
                value={emailDomain}
                onChange={selectDomain}
              >
                <option value="직접입력">직접입력</option>
                <option value="naver.com" selected>
                  naver.com
                </option>
                <option value="hanmail.net">hanmail.net</option>
                <option value="hotmail.com">hotmail.com</option>
                <option value="nate.com">nate.com</option>
                <option value="yahoo.co.kr">yahoo.co.kr</option>
                <option value="empas.com">empas.com</option>
                <option value="dreamwiz.com">dreamwiz.com</option>
                <option value="freechal.com">freechal.com</option>
                <option value="lycos.co.kr">lycos.co.kr</option>
                <option value="korea.com">korea.com</option>
                <option value="gmail.com">gmail.com</option>
                <option value="hanmir.com">hanmir.com</option>
                <option value="paran.com">paran.com</option>
              </select>
            </div>
            <div className={styles.errormessage}>
              {userIDError
                ? userIDError
                : passwordError
                ? passwordError
                : passwordCheckError
                ? passwordCheckError
                : nicknameError
                ? nicknameError
                : ""}
            </div>
            <button
              onClick={submit}
              disabled={
                userID === "" ||
                password === "" ||
                passwordChk === "" ||
                nickname === "" ||
                emailID === "" ||
                emailDomain === "" ||
                userIDError !== "" ||
                passwordError !== "" ||
                passwordCheckError !== "" ||
                nicknameError !== "" ||
                userIdDuplicationText !== "사용가능" ||
                nicknameDuplicationText !== "사용가능"
              }
              className={styles.signup}
            >
              Sign up
            </button>
          </form>
        </div>
      )}
    </div>
  );
}

export default SignUp;
