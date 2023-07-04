import { useState } from "react";
import styles from "./Login.module.css";
import { BsXLg } from "react-icons/bs";

function Login({ message, popupOpen }) {
  const [ID, setID] = useState("");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState("");

  const onIdChange = (e) => {
    setID(e.target.value);
  };
  const onPwChange = (e) => {
    setPassword(e.target.value);
  };

  const login = (e) => {};

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
        <h5 className={styles.message}>{message}</h5>
        <form className={styles.boxinput} onSubmit={login}>
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
        </form>
      </div>
    </div>
  );
}

export default Login;
