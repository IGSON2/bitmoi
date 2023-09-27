import styles from "./H_NavBar.module.css";
import { Link } from "react-router-dom";
import ProfileBox from "../ProfileBox/ProfileBox";
import { useEffect, useState } from "react";
import checkAccessTokenValidity from "../backendConn/checkAccessTokenValidity";

function H_NavBar() {
  const menubutton = ["HOME", "COMPETITION", "PRACTICE", "RANK", "AD BIDDING"];
  const [userInfo, setUserinfo] = useState();

  useEffect(() => {
    const verifyToken = async () => {
      const userInfo = await checkAccessTokenValidity();
      if (userInfo) {
        setUserinfo(userInfo);
      }
    };
    verifyToken();
  }, []);

  return (
    <div className={styles.navbar}>
      {menubutton.map((menu, idx) => {
        if (menu === menubutton[0]) {
          return (
            <Link key={idx} className={styles.navmenu} to={`/`}>
              {menu}
            </Link>
          );
        } else {
          return (
            <Link
              key={idx}
              className={styles.navmenu}
              to={`/${menu.replace(" ", "-").toLowerCase()}`}
            >
              {menu}
            </Link>
          );
        }
      })}
      <div className={styles.profilediv}>
        <ProfileBox userInfo={userInfo} />
      </div>
    </div>
  );
}

export default H_NavBar;
