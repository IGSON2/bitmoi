import { useEffect, useState } from "react";
import styles from "./Login.module.css";
import { BsXLg } from "react-icons/bs";
import { Link } from "react-router-dom";
import axiosClient from "../../component/backendConn/axiosClient";

function Login() {
  const [ID, setID] = useState("");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState("");
  const onIdChange = (e) => {
    setID(e.target.value);
  };
  const onPwChange = (e) => {
    setPassword(e.target.value);
  };

  const login = (e) => {
    e.preventDefault(e);
    try {
      const response = axiosClient.post("/user/login", {
        user_id: ID,
        password: password,
      });
      console.log(response);
    } catch (error) {
      console.error(error);
    }

    // loginPromise
    //   .then((r) => {
    //     const res = r.json();
    //     console.log(res);
    //     localStorage.setItem("access_token", res.access_token);
    //     localStorage.setItem("access_token", res.refresh_token);
    //     setErrorMsg("");
    //     setIsLogined(true);
    //   })
    //   .catch((error) => {
    //     setErrorMsg(error);
    //   });
  };
  return (
    <div className={styles.loginwindow}>
      <div className={styles.bg}></div>
      <div className={styles.popupbody}>
        <div className={styles.closebutton}>
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
        <Link to={"/signup"}>
          <button className={styles.signup}>Sign up</button>
        </Link>
      </div>
    </div>
  );
}

export default Login;
