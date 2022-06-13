import { Link } from "react-router-dom";
import styles from "./V_NavBar.module.css";

function V_Navbar() {
  const menubutton = [
    "Home",
    "Competition",
    "Practice",
    "Ranking",
    // "Community",
  ];
  return (
    <div className={styles.navpage}>
      <div className={styles.navbar}>
        {menubutton.map((menu, idx) => {
          if (menu === "Home") {
            return (
              <Link key={idx} className={styles.navmenu} to={`/`}>
                <button className={styles.navmenubutton}>{menu}</button>
              </Link>
            );
          } else {
            return (
              <Link key={idx} className={styles.navmenu} to={`/${menu}`}>
                <button className={styles.navmenubutton}>{menu}</button>
              </Link>
            );
          }
        })}
      </div>
    </div>
  );
}

export default V_Navbar;
