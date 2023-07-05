import { useState } from "react";
import styles from "./Login.module.css";
import { BsXLg } from "react-icons/bs";
import PostLogin from "../backendConn/PostLogin";
import { SaveAccessToken, SaveRefreshToken } from "../Token/Token";

function Login({ popupOpen, setUserInfo, setIsLogined }) {
  const [ID, setID] = useState("");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState("");

  const onIdChange = (e) => {
    setID(e.target.value);
  };
  const onPwChange = (e) => {
    setPassword(e.target.value);
  };

  const login = () => {
    const loginPromise = PostLogin(
      "http://bitmoi.co.kr:5000/user/login",
      ID,
      password
    );
    loginPromise
      .then((res) => {
        SaveAccessToken(res.access_token);
        SaveRefreshToken(res.refresh_token);
        setUserInfo(res.user);
        setErrorMsg("");
        setIsLogined(true);
      })
      .catch((error) => {
        setErrorMsg(error);
      });
  };

  const closePopup = () => {
    popupOpen(false);
  };

  return (
    <div className={styles.loginwindow}>
      <div className={styles.bg} onClick={closePopup}></div>
      <div className={styles.popupbody}>
        <div className={styles.closebutton} onClick={closePopup}>
          <BsXLg />
        </div>
        <h1 className={styles.title}>로그인</h1>
        <h5 className={styles.message}>
          비트모이에 로그인하여 경쟁에 참여해 보세요.
        </h5>
        <form className={styles.inputform} onSubmit={login}>
          <input
            className={styles.box}
            onChange={onIdChange}
            value={ID}
            placeholder="ID"
          />
          <input
            className={styles.box}
            onChange={onPwChange}
            type="password"
            value={password}
            placeholder="password"
          />
          <button className={styles.login} onClick={login}>
            Login
          </button>
        </form>
        <p>{errorMsg}</p>
      </div>
    </div>
  );
}

export default Login;
