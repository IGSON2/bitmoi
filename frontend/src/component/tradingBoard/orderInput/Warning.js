import styles from "./warning.module.css";
function Warning({
  profitWarning,
  lossWarning,
  levWarning,
  quanWarning,
  loginWarning,
}) {
  var warningTxt = "";
  if (loginWarning !== "") {
    warningTxt = loginWarning;
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

  return (
    <div className={styles.warningdiv}>
      <p>{warningTxt}</p>
    </div>
  );
}
export default Warning;
