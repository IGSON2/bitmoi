import { useEffect, useState } from "react";
import Login from "../login/Login";
import { BsBoxArrowRight } from "react-icons/bs";
import { IoIosPerson } from "react-icons/io";
import logo from "../images/logosmall.png";
import styles from "./ProfileBox.module.css";
import { Link } from "react-router-dom";

function ProfileBox() {
  const [loginClick, setLoginClick] = useState(false);
  const [isLogined, setIsLogined] = useState(false);
  const [authLoad, setAuthLoad] = useState(false);

  const loginPopup = () => {
    setLoginClick(true);
  };

  const [profileURL, setProfileURL] = useState("");
  const [loginTxt, setLoginTxt] = useState("");
  const [openProfile, setOpenProfile] = useState(false);
  const profileClick = () => {
    if (authLoad) {
      setOpenProfile((current) => !current);
    }
  };

  useEffect(() => {
    setOpenProfile(false);
    setLoginTxt("Loading...");
  }, []);

  const logOut = async () => {};

  return (
    <div className={styles.profiebox}>
      {isLogined ? (
        openProfile ? (
          <div className={styles.openedProfile}>
            <img
              className={styles.profileImg}
              src={profileURL ? profileURL : logo}
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
            src={profileURL ? profileURL : logo}
            onClick={profileClick}
          ></img>
        )
      ) : (
        <button
          className={styles.loginbutton}
          onClick={loginPopup}
          disabled={!authLoad}
        >
          {loginTxt}
        </button>
      )}
      {loginClick ? (
        <Login
          message={"비트모이에 로그인하여 트레이딩 전적을 저장해 보세요."}
          popupOpen={setLoginClick}
        />
      ) : null}
    </div>
  );
}

export default ProfileBox;
