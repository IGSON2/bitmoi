import { useRef, useState } from "react";
import styles from "./Login.module.css";
import { Link } from "react-router-dom";
import axiosClient from "../../component/backendConn/axiosClient";
import H_NavBar from "../../component/navbar/H_NavBar";

function Login() {
  const [ID, setID] = useState("");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState("");
  const inputButtonRef = useRef(null);
  const onIdChange = (e) => {
    setID(e.target.value);
  };
  const onPwChange = (e) => {
    setPassword(e.target.value);
  };

  const login = async (e) => {
    e.preventDefault(e);
    try {
      const response = await axiosClient.post("/basic/user/login", {
        user_id: ID,
        password: password,
      });
      localStorage.setItem("accessToken", response.data.access_token);
      localStorage.setItem("refreshToken", response.data.refresh_token);
      setErrorMsg("");
      window.location.replace("/competition");
    } catch (error) {
      console.error(error);
      setErrorMsg("ID 또는 PW를 확인해 주세요");
    }
  };

  const pressEnter = () => {};

  return (
    <div className={styles.loginwindow}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      <div className={styles.popupbody}>
        <p className={styles.subtitle}>바로 시작하는 모의투자</p>
        <h1 className={styles.title}>BITMOI</h1>
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
          <button style={{ display: "none" }}></button>
          {errorMsg ? <p className={styles.errormessage}>{errorMsg}</p> : null}
        </form>
        <button className={styles.login} onClick={login}>
          Login
        </button>
        <div className={styles.signup}>
          <Link to={"/signup"}>Sign up</Link>
        </div>
      </div>
    </div>
  );
}

export default Login;
