import styles from "./warning.module.css";
function Warning({
  profitWarning,
  lossWarning,
  levWarning,
  quanWarning,
  tokenWarning,
}) {
  var warningTxt = "";
  if (tokenWarning !== "") {
    warningTxt = tokenWarning;
  } else {
    if (profitWarning !== "") {
      warningTxt = profitWarning;
    } else {
      if (lossWarning !== "") {
        warningTxt = lossWarning;
      } else {
        if (levWarning !== "") {
          warningTxt = levWarning;
        } else {
          if (quanWarning !== "") {
            warningTxt = quanWarning;
          } else {
            warningTxt = "";
          }
        }
      }
    }
  }

  const goFreetoken = () => {
    window.location.replace("/freetoken");
  };

  return (
    <div className={styles.warningdiv}>
      <p>{warningTxt}</p>
      {tokenWarning === "도전에 사용할 MOI 토큰이 부족합니다." ? (
        <a className={styles.gofreetoken} href="/freetoken" target="_blank">
          무료 토큰 받기
        </a>
      ) : null}
    </div>
  );
}
export default Warning;
