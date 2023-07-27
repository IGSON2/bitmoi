import { useState } from "react";
import { BsBoxArrowRight } from "react-icons/bs";
import { IoIosPerson } from "react-icons/io";
import logo from "../images/logosmall.png";
import styles from "./ProfileBox.module.css";
import { Link } from "react-router-dom";

function ProfileBox(props) {
  const routeLogin = () => {
    window.location.href = "/login";
  };

  const [openProfile, setOpenProfile] = useState(false);
  const profileClick = () => {
    setOpenProfile((current) => !current);
  };

  const logOut = () => {
    localStorage.removeItem("accessToken");
  };

  console.log(props.userInfo);

  return (
    <div className={styles.profiebox}>
      {props.userInfo ? (
        openProfile ? (
          <div className={styles.openedProfile}>
            <img
              className={styles.profileImg}
              src={props.userInfo.photo_url ? props.userInfo.photo_url : logo}
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
            src={props.userInfo.photo_url ? props.userInfo.photo_url : logo}
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
