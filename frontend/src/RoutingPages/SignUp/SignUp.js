import { useEffect, useState } from "react";
import styles from "./SignUp.module.css";
import H_NavBar from "../../component/navbar/H_NavBar";
function SignUp() {
  const [passwordChk, setPasswordChk] = useState("");
  const [isLogined, setIsLogined] = useState(false);
  const onSubmit = (e) => {
    e.preventDefault();
  };
  useEffect(() => {
    onAuthStateChanged(auth, (user) => {
      if (user) {
        setIsLogined(true);
      } else {
        setIsLogined(false);
      }
    });
  }, []);
  const pwChkChange = (e) => {};

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
              <label htmlFor="email">이메일</label>
              <input id="email"></input>
            </div>
            <div className={styles.field}>
              <label htmlFor="pw">비밀번호</label>
              <input id="pw" type="password"></input>
            </div>
            <div className={styles.field}>
              <label htmlFor="pwcheck">비밀번호 확인</label>
              <input
                id="pwcheck"
                type="password"
                value={passwordChk}
                onChange={pwChkChange}
              ></input>
            </div>
            <div className={styles.field}>
              <label htmlFor="nickname">닉네임</label>
              <input id="nickname"></input>
            </div>
            <button className={styles.signup}>Sign up</button>
            <a className={styles.back} href="/practice">
              가입없이 연습모드만 먼저 해볼게요!
            </a>
          </form>
        </div>
      )}
    </div>
  );
}

export default SignUp;
