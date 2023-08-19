import { Link } from "react-router-dom";
import styles from "./V_NavBar.module.css";

function V_Navbar() {
  const menubutton = ["HOME", "COMPETITION", "PRACTICE", "RANK", "AD BIDDING"];
  return (
    <div className={styles.navpage}>
      <div className={styles.navbar}>
        {menubutton.map((menu, idx) => {
          if (menu === menubutton[0]) {
            return (
              <Link key={idx} className={styles.navmenu} to={`/`}>
                <button className={styles.navmenubutton}>{menu}</button>
              </Link>
            );
          } else {
            return (
              <Link
                key={idx}
                className={styles.navmenu}
                to={`/${menu.replace(" ", "-").toLowerCase()}`}
              >
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
