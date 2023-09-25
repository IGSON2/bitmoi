import styles from "./Countbox.module.css";

function Countbox({ number, unit }) {
  return (
    <div className={styles.wrapper}>
      <div className={styles.number}>{number}</div>
      {/* <div className={styles.unit}>{unit}</div> */}
    </div>
  );
}

export default Countbox;
