import { useEffect, useState } from "react";
import styles from "./SignUp.module.css";
import H_NavBar from "../../component/navbar/H_NavBar";
import { LoadAccessToken } from "../../component/Token/Token";
function SignUp() {
  const [userID, setUserID] = useState("");
  const [userIDCheckError, setUserIDCheckError] = useState("");
  const [emailID, setEmailID] = useState("");
  const [emailDomain, setEmailDomain] = useState("");
  const [selectDomainDisable, setSelectDomainDisable] = useState(true);
  const [password, setPassword] = useState("");
  const [passwordChk, setPasswordChk] = useState("");
  const [nickname, setNickname] = useState("");
  const [isLogined, setIsLogined] = useState(false);
  const [nicknameCheckError, setNicknameCheckError] = useState("");

  const onSubmit = (e) => {
    e.preventDefault();
  };

  const userIDChange = (e) => {
    setUserID(e.target.value);
  };

  const userIDCheck = () => {
    fetch("http://bitmoi.co.kr:5000/user/checkid?user_id=" + userID)
      .then((res) => {
        if (res.ok) {
          setUserIDCheckError("");
          return;
        }
      })
      .catch((error) => {
        setUserIDCheckError(error);
      });
  };

  const emailIDChange = (e) => {
    setEmailID(e.target.value + "@");
  };

  const selectDomain = (e) => {
    if (e.target.value === "1") {
      setSelectDomainDisable(false);
    } else {
      setEmailDomain(e.target.value);
      setSelectDomainDisable(true);
    }
  };
  const typingDomain = (e) => {
    setEmailDomain(e.target.value);
  };

  const pwChange = (e) => {
    setPassword(e.target.value);
  };

  const pwChkChange = (e) => {
    setPasswordChk(e.target.value);
  };

  const nicknameChange = (e) => {
    setNickname(e.target.value);
  };

  const nicknameCheck = () => {
    fetch("http://bitmoi.co.kr:5000/user/checknickname?nickname=" + nickname)
      .then((res) => {
        if (res.ok) {
          setNicknameCheckError("");
          return;
        }
      })
      .catch((error) => {
        setNicknameCheckError(error);
      });
  };

  useEffect(() => {
    if (LoadAccessToken() !== "") {
      setIsLogined(true);
    }
  }, []);

  return (
    <div className={styles.signupdiv}>
      <H_NavBar></H_NavBar>
      {isLogined ? (
        <div className={styles.warning}>
          <h1>잘못된 접근입니다!</h1>
          <a href="/">BACK</a>
        </div>
      ) : (
        <div className={styles.formdiv}>
          <form className={styles.forms} onSubmit={onSubmit}>
            <h3 className={styles.welcome}>
              시뮬레이션 모의투자 비트모이에 오신 걸 환영합니다!
            </h3>

            <div className={styles.field}>
              <label htmlFor="id">아이디</label>
              <input
                id="id"
                type="text"
                placeholder="ID"
                pattern="^[a-zA-Z0-9]{5,15}$"
                value={userID}
                onChange={userIDChange}
              ></input>
              <button onClick={userIDCheck}>중복확인</button>
            </div>
            {userIDCheckError ? <div>{userIDCheckError}</div> : null}

            <div className={styles.field}>
              <label htmlFor="pw">비밀번호</label>
              <input
                id="pw"
                type="password"
                placeholder="password"
                pattern="^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,16}$"
                value={password}
                onChange={pwChange}
              ></input>
            </div>

            <div className={styles.field}>
              <label htmlFor="pwcheck">비밀번호 확인</label>
              <input
                id="pwcheck"
                type="password"
                placeholder="password"
                pattern="^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,16}$"
                value={passwordChk}
                onChange={pwChkChange}
              ></input>
            </div>

            <div className={styles.field}>
              <label htmlFor="nickname">닉네임</label>
              <input
                id="nickname"
                type="text"
                value={nickname}
                onChange={nicknameChange}
              ></input>
              <button onClick={nicknameCheck}>중복확인</button>
            </div>
            {nicknameCheckError ? <div>{nicknameCheckError}</div> : null}

            <div className={styles.field}>
              <label htmlFor="emailID">이메일</label>
              <input
                id="emailID"
                value={emailID}
                onChange={emailIDChange}
              ></input>
              <input
                disabled={selectDomainDisable}
                value={emailDomain}
                onChange={typingDomain}
              ></input>
              <select
                id="selectEmailDomain"
                value={emailDomain}
                onChange={selectDomain}
              >
                <option value="1">직접입력</option>
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

            <button className={styles.signup}>Sign up</button>
          </form>
        </div>
      )}
    </div>
  );
}

export default SignUp;
