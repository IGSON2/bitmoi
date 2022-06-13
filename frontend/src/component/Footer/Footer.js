import styles from "./Footer.module.css";
import { BsInstagram, BsFacebook, BsTwitter } from "react-icons/bs";
function Footer() {
  return (
    <div className={styles.footer}>
      <div className={`${styles.insta} ${styles.icon}`}>
        <BsInstagram />
      </div>
      <div className={`${styles.facebook} ${styles.icon}`}>
        <BsFacebook />
      </div>
      <div className={`${styles.twitter} ${styles.icon}`}>
        <BsTwitter />
      </div>
      <div className={styles.nav}>
        <a href="/">Home</a>
        <a href="/practice">Practice</a>
        <a href="/competition">Competition</a>
      </div>
      <div className={styles.copyright}>Copyright &copy; 2022 IGSON</div>
    </div>
  );
}

export default Footer;
