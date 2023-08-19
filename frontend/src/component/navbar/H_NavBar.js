import styles from "./H_NavBar.module.css";
import { Link } from "react-router-dom";

function H_NavBar() {
  const menubutton = ["HOME", "COMPETITION", "PRACTICE", "RANK", "AD BIDDING"];
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
    </div>
  );
}

export default H_NavBar;
