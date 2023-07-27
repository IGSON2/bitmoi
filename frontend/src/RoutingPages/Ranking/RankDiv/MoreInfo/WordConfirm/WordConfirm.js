import styles from "./WordConfirm.module.css";
function WordConfirm({ popupOpen, comment }) {
  const closePopup = () => {
    popupOpen(false);
  };
  const yesnoclick = (e) => {
    switch (e.target.innerText) {
      case "네":
        fetch("http://bitmoi.co.kr:5000/moreinfo", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            comment: comment,
            user: "temporary_user",
          }),
        }).then(window.location.reload());
        break;
      case "아니오":
        popupOpen(false);
        break;
    }
  };

  return (
    <div className={styles.confirmwindow}>
      <div className={styles.bg} onClick={closePopup}></div>
      <div className={styles.popupbody}>
        <h3>소감은 최초 한 번만 등록 가능합니다. 이대로 등록 할까요?</h3>
        <div className={styles.yesno}>
          <button onClick={yesnoclick}>네</button>
          <button onClick={yesnoclick}>아니오</button>
        </div>
      </div>
    </div>
  );
}

export default WordConfirm;
