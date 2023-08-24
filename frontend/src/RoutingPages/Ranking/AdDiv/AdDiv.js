import styles from "./AdDiv.module.css";
import mockup from "../../../component/images/mockup_rank.png";

function AdDiv() {
  const adClick = () => {
    window.location.replace("/ad-bidding/rank");
  };
  return (
    <div className={styles.addiv} onClick={adClick}>
      <img className={styles.adimage} src={mockup}></img>
    </div>
  );
}

export default AdDiv;
