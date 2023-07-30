import { useState } from "react";
import logo from "../images/logosmall.png";
import styles from "./ProfileBox.module.css";
import HorizontalLine from "../lines/HorizontalLine";
import { BsXLg } from "react-icons/bs";

function ProfileBox(props) {
  const routeLogin = () => {
    window.location.href = "/login";
  };

  const [openProfile, setOpenProfile] = useState(false);

  const profileClick = () => {
    setOpenProfile((current) => !current);
  };

  const closePopup = () => {
    setOpenProfile(false);
  };

  const logOut = () => {
    localStorage.removeItem("accessToken");
  };
  return (
    <div className={styles.profiebox}>
      {props.userInfo ? (
        <img
          className={styles.profileImg}
          src={props.userInfo.photo_url ? props.userInfo.photo_url : logo}
          onClick={profileClick}
        ></img>
      ) : (
        <button className={styles.loginbutton} onClick={routeLogin}>
          login
        </button>
      )}
      {openProfile ? (
        <div className={styles.userinfo}>
          <div className={styles.top}>
            <div className={styles.namebox}>
              <span className={styles.name}>{props.userInfo.nickname}</span>
              <span className={styles.welcome}> 님 안녕하세요!</span>
            </div>
            <div className={styles.closebutton} onClick={closePopup}>
              <BsXLg />
            </div>
          </div>
          <HorizontalLine />
          <div className={styles.middle}></div>
        </div>
      ) : null}
    </div>
  );
}

export default ProfileBox;
