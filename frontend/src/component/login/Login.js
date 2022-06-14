import {
  browserLocalPersistence,
  browserSessionPersistence,
  getAuth,
  GoogleAuthProvider,
  setPersistence,
  signInWithEmailAndPassword,
  signInWithPopup,
  signOut,
} from "firebase/auth";
import { useState } from "react";
import styles from "./Login.module.css";
import { BsXLg } from "react-icons/bs";
import googleLogin from "../images/btn_google_signin_light_normal_web@2x.png";
import googleFocused from "../images/btn_google_signin_light_focus_web@2x.png";
function Login({ message, popupOpen }) {
  const auth = getAuth();
  const provider = new GoogleAuthProvider();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [warningMsg, setWarningMsg] = useState("");
  const [rememberme, setRememberme] = useState(false);
  const [googleButton, setGoogleButton] = useState(googleLogin);
  const onChange = (e) => {
    switch (e.target.placeholder) {
      case "Email":
        setEmail(e.target.value);
        break;
      case "Password":
        setPassword(e.target.value);
        break;
    }
  };
  const checkboxChange = (e) => {
    if (e.target.checked) {
      setRememberme(true);
    } else {
      setRememberme(false);
    }
  };
  const signIn = async (e) => {
    e.preventDefault();

    try {
      console.log(auth);
      await signInWithEmailAndPassword(auth, email, password);
      setWarningMsg("");
      popupOpen(false);
    } catch (err) {
      console.log(err.code);
      switch (err.code) {
        case "auth/invalid-email":
          setWarningMsg("유효하지 않은 Email 주소 입니다.");
          break;
        case "auth/user-not-found":
          setWarningMsg("Email 또는 비밀번호가 올바른지 확인해 주세요.");
          break;

        case "auth/wrong-password":
          setWarningMsg("Email 또는 비밀번호가 올바른지 확인해 주세요.");
          break;

        case "auth/too-many-requests":
          setWarningMsg("잦은 요청으로 잠시후에 다시 시도가 가능해 집니다.");
          break;
      }
    }
  };

  const closePopup = () => {
    popupOpen(false);
  };

  const gLogion = () => {
    if (rememberme) {
      setPersistence(auth, browserLocalPersistence);
    } else {
      setPersistence(auth, browserSessionPersistence);
    }
    signInWithPopup(auth, provider).then(popupOpen(false));
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
        <div className={styles.boxinput}>
          <input
            className={styles.box}
            id="box"
            type="checkbox"
            onChange={checkboxChange}
          />
          <label className={styles.boxlabel} htmlFor="box">
            Remember me
          </label>
        </div>
        <div className={styles.googlelogindiv}>
          <img
            className={styles.googleloginbutton}
            src={googleButton}
            onMouseOver={() => {
              setGoogleButton(googleFocused);
            }}
            onMouseLeave={() => {
              setGoogleButton(googleLogin);
            }}
            onClick={gLogion}
          />
        </div>
      </div>
    </div>
  );
}

export default Login;
