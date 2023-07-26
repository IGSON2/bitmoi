import styles from "./ChartHeader.module.css";
import horizontalImg from "../../images/horizontal.png";
import rulerImg from "../../images/ruler.png";
import eraserImg from "../../images/eraser.png";
import keyboardImg from "../../images/keyboard.png";
import mouseImg from "../../images/mouse.png";
import markImg from "../../images/marker.png";
import ProfileBox from "../../ProfileBox/ProfileBox";

function ChartHeader(props) {
  const fiveMinute = () => {
    props.getChartData("5m");
  };
  const fifteenMinute = () => {
    props.getChartData("15m");
  };
  const oneHour = () => {
    props.getChartData("1h");
  };
  const fourHour = () => {
    props.getChartData("4h");
  };
  const markButtonClick = () => {
    props.setToolBar("mark");
    props.setActive("mark");
    if (props.active === "mark") {
      props.setToolBar("NonSelected");
      props.setActive("");
    }
  };
  const horizontalButtonClick = () => {
    props.setToolBar("horizon");
    props.setActive("horizon");
    if (props.active === "horizon") {
      props.setToolBar("NonSelected");
      props.setActive("");
    }
  };
  const rulerButtonClick = () => {
    props.setToolBar("ruler");
    props.setActive("ruler");
    if (props.active === "ruler") {
      props.setToolBar("NonSelected");
      props.setActive("");
    }
  };
  const eraserButtonClick = () => {
    props.setToolBar("escape");
    props.setActive("ruler");
  };

  const modechange = () => {
    if (props.toolBar === "NonSelected" || props.toolBar === "") {
      props.setToolBar("ruler");
      props.setActive("ruler");
    } else {
      props.setToolBar("NonSelected");
      props.setActive("");
    }
  };

  const helppopup = () => {
    alert(
      `Hotkey Manual\n=================\nRuler : Shift + Click\nMakrer : Ctrl + Click\nPriceLine : Alt + Click\nEraseAll : Click Chart`
    );
  };

  return (
    <div className={styles.ChartHeader}>
      <div className={styles.title}>
        <div className={styles.name}>{props.name}</div>
      </div>
      <div className={styles.addinfo}>
        <div className={styles.entrytime} title={"시뮬레이션 시점"}>
          <div className={styles.childname}>Entrytime</div>
          <div className={styles.childvalue}>{props.entryTime}</div>
        </div>
        <div
          className={styles.btcratio}
          title={"최근 한 달간 비트코인의 거래량 대비"}
        >
          <div className={styles.childname}>Lightness</div>
          <div className={styles.childvalue}>{props.btcRatio.toFixed(3)}%</div>
        </div>
      </div>
      <div className={styles.buttons}>
        <div className={styles.buttonsparent}>
          <button
            className={styles.buttonschildren}
            style={
              props.headerInterval === "5m"
                ? {
                    color: "#4f62cc",
                    fontWeight: 700,
                  }
                : { color: "#555860" }
            }
            disabled={props.headerInterval === "submit" ? true : false}
            onClick={fiveMinute}
          >
            5M
          </button>
          <button
            className={styles.buttonschildren}
            style={
              props.headerInterval === "15m"
                ? {
                    color: "#4f62cc",
                    fontWeight: 700,
                  }
                : { color: "#555860" }
            }
            disabled={props.headerInterval === "submit" ? true : false}
            onClick={fifteenMinute}
          >
            15M
          </button>
          <button
            className={styles.buttonschildren}
            style={
              props.headerInterval === "1h" || props.headerInterval === "submit"
                ? {
                    color: "#4f62cc",
                    fontWeight: 700,
                  }
                : { color: "#555860" }
            }
            disabled={props.headerInterval === "submit" ? true : false}
            onClick={oneHour}
          >
            1H
          </button>
          <button
            className={styles.buttonschildren}
            style={
              props.headerInterval === "4h"
                ? {
                    color: "#4f62cc",
                    fontWeight: 700,
                  }
                : { color: "#555860" }
            }
            disabled={props.headerInterval === "submit" ? true : false}
            onClick={fourHour}
          >
            4H
          </button>
        </div>
        <div className={styles.buttonsparent}>
          <button className={styles.imgbox} onClick={modechange}>
            {props.toolBar === "NonSelected" || props.toolBar === "" ? (
              <img
                alt="Hotkey Mode"
                title="Hotkey Mode"
                src={keyboardImg}
              ></img>
            ) : (
              <img alt="Mouse Mode" title="Mouse mode" src={mouseImg}></img>
            )}
          </button>
          <button
            className={styles.imgbox}
            style={
              props.active === "mark"
                ? { backgroundColor: "rgba(79, 98, 197,0.6)" }
                : null
            }
            onClick={markButtonClick}
          >
            <img alt="Marker" src={markImg}></img>
          </button>
          <button
            className={styles.imgbox}
            style={
              props.active === "horizon"
                ? { backgroundColor: "rgba(79, 98, 197,0.6)" }
                : null
            }
            onClick={horizontalButtonClick}
          >
            <img alt="Price Line" src={horizontalImg}></img>
          </button>
          <button
            className={styles.imgbox}
            style={
              props.active === "ruler"
                ? { backgroundColor: "rgba(79, 98, 197,0.6)" }
                : null
            }
            onClick={rulerButtonClick}
          >
            <img alt="Ruler" src={rulerImg}></img>
          </button>
          <button
            className={styles.imgbox}
            style={
              props.active === "escape"
                ? { backgroundColor: "rgba(79, 98, 197,0.6)" }
                : null
            }
            onClick={eraserButtonClick}
          >
            <img alt="Eraser" src={eraserImg}></img>
          </button>
        </div>
      </div>
      {props.toolBar === "NonSelected" ? (
        <button className={styles.helpbutton} onClick={helppopup}>
          ?
        </button>
      ) : (
        <div className={styles.blank}></div>
      )}
      <div className={styles.profilediv}>
        <ProfileBox userInfo={props.userInfo} />
      </div>
    </div>
  );
}

export default ChartHeader;
