import styles from "./H_NavBar.module.css";
import { Link } from "react-router-dom";
import ProfileBox from "../ProfileBox/ProfileBox";

function H_NavBar() {
  const menubutton = [
    "HOME",
    "COMPETITION",
    "PRACTICE",
    "RANK",
    // "Community",
  ];
  return (
    <div className={styles.navbar}>
      {menubutton.map((menu, idx) => {
        if (menu === "HOME") {
          return (
            <Link key={idx} className={styles.navmenu} to={`/`}>
              {menu}
            </Link>
          );
        } else {
          return (
            <Link key={idx} className={styles.navmenu} to={`/${menu}`}>
              {menu}
            </Link>
          );
        }
      })}
    </div>
  );
}

export default H_NavBar;
