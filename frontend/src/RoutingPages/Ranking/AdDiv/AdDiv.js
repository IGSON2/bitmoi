import styles from "./AdDiv.module.css";
import { useEffect } from "react";

function AdDiv({ imgLink }) {
  const adClick = () => {
    window.open("/ad-bidding/rank", "_blank");
  };
  return (
    <div className={styles.addiv} onClick={adClick}>
      <img className={styles.adimage} src={imgLink}></img>
    </div>
  );
}

export default AdDiv;
