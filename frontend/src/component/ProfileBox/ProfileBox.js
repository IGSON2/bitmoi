import { useState } from "react";
import logo from "../images/logo.png";
import mypage from "../images/my_page.png";
import logout from "../images/logout.png";
import styles from "./ProfileBox.module.css";
import HorizontalLine from "../lines/HorizontalLine";
import VerticalLine from "../lines/VerticalLine";
import { BsXLg } from "react-icons/bs";
import Wallet from "../Wallet/Wallet";

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

  const gotoMypage = () => {
    window.location.href = "/mypage";
  };
  const logOut = () => {
    localStorage.removeItem("accessToken");
    localStorage.removeItem("refreshToken");
    window.location.reload();
  };
  return (
    <div className={styles.profiebox}>
      {props.userInfo && props.userInfo.user_id !== "" ? (
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
          <div className={styles.middle}>
            <div className={styles.moitoken}>
              <Wallet />
            </div>
            <VerticalLine />
            <div className={styles.mypage}>
              <img src={mypage} onClick={gotoMypage}></img>
            </div>
            <div className={styles.logout}>
              <img src={logout} onClick={logOut}></img>
            </div>
          </div>
        </div>
      ) : null}
    </div>
  );
}

export default ProfileBox;
