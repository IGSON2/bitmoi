import styles from "./Header.module.css";
function Header() {
  return (
    <div className={styles.header}>
      <div className={styles.fieldshort}>번호</div>
      <div className={styles.fieldlong}>시간</div>
      <div className={styles.fieldlong}>페어이름</div>
      <div className={styles.fieldlong}>복기지점</div>
      <div className={styles.fieldshort}>스테이지</div>
      <div className={styles.field}>포지션</div>
      <div className={styles.fieldshort}>레버리지</div>
      <div className={styles.field}>진입가격</div>
      <div className={styles.field}>정리가격</div>
      <div className={styles.fieldshort}>보유시간</div>
      <div className={styles.field}>수익금</div>
      <div className={styles.field}>수익률</div>
    </div>
  );
}

export default Header;
