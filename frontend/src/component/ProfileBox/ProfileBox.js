import { useEffect, useState } from "react";
import { BsBoxArrowRight } from "react-icons/bs";
import { IoIosPerson } from "react-icons/io";
import logo from "../images/logosmall.png";
import styles from "./ProfileBox.module.css";
import { Link } from "react-router-dom";
import checkAccessTokenValidity from "../backendConn/checkAccessTokenValidity";

function ProfileBox() {
  const [isLogined, setIsLogined] = useState(false);
  const routeLogin = () => {
    window.location.href = "/login";
  };
  const [userInfo, setUserInfo] = useState({
    user_id: "",
    nickname: "",
    email: "",
    photo_url: "",
  });
  const [openProfile, setOpenProfile] = useState(false);
  const profileClick = () => {
    setOpenProfile((current) => !current);
  };

  const logOut = () => {
    localStorage.removeItem("accessToken");
  };

  return (
    <div className={styles.profiebox}>
      {isLogined ? (
        openProfile ? (
          <div className={styles.openedProfile}>
            <img
              className={styles.profileImg}
              src={userInfo.photo_url ? userInfo.photo_url : logo}
              onClick={profileClick}
            ></img>
            <div className={styles.nameoptions}>
              <div className={styles.username}></div>
              <div className={styles.options}>
                <Link title="My Page" to={"/myscore"}>
                  <IoIosPerson />
                </Link>
                <div title="Log out" onClick={logOut}>
                  <BsBoxArrowRight />
                </div>
              </div>
            </div>
          </div>
        ) : (
          <img
            className={styles.profileImg}
            src={userInfo.photo_url ? userInfo.photo_url : logo}
            onClick={profileClick}
          ></img>
        )
      ) : (
        <button className={styles.loginbutton} onClick={routeLogin}>
          login
        </button>
      )}
    </div>
  );
}

export default ProfileBox;
