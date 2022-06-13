import styles from "./Mobile.module.css";
function Mobile() {
  return (
    <div className={styles.confirmwindow}>
      <div className={styles.bg}></div>
      <div className={styles.popupbody}>
        <h4>현재 모바일 버전은 준비중에 있습니다.</h4>
      </div>
    </div>
  );
}

export default Mobile;
