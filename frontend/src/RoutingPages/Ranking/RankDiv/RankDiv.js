import styles from "./RankDiv.module.css";
import { BsChatQuote } from "react-icons/bs";
import { useEffect, useState } from "react";
import MoreInfo from "./MoreInfo/MoreInfo";
import axiosClient from "../../../component/backendConn/axiosClient";

function RankDiv({ index, obj }) {
  const [moreInfo, setMoreInfo] = useState(false);
  const [certified, setCertified] = useState(false);
  const colorchanger = (num) => {
    var color = "";
    switch (num) {
      case 1:
        color = "gold";
        break;
      case 2:
        color = "silver";
        break;
      case 3:
        color = "#E69567";
        break;
      default:
        color = "black";
    }
    return color;
  };
  const [data, setData] = useState({
    comment: "",
    scoreid: "",
    avglev: 0,
    avgpnl: 0,
    avgroe: 0,
    stagearray: [
      {
        name: "",
        date: "",
        roe: 0,
      },
    ],
  });
  const getMoreInfo = () => {
    if (!moreInfo) {
      const response = axiosClient.get(
        `http://bitmoi.co.kr:5000/moreinfo?userid=${obj.user_id}&scoreid=${obj.score_id}`
      );
      console.log(response.data);
      setData(response.data);
    }

    setMoreInfo((current) => !current);
  };

  return (
    <div className={styles.userdiv}>
      <div className={styles.onlyinfo}>
        <div
          className={`${styles.no} ${styles.field}`}
          style={{ color: `${colorchanger(index)}` }}
        >
          {index}
        </div>
        <div className={`${styles.pic} ${styles.field}`}>
          <img className={styles.photo} src={obj.photo_url} />
        </div>
        <div className={`${styles.name} ${styles.field}`}>{obj.nickname}</div>
        <div className={`${styles.score}  ${styles.field}`}>
          {obj.final_balance}
        </div>
        <button className={styles.openbutton} onClick={getMoreInfo}>
          <BsChatQuote />
        </button>
      </div>
      <div className={styles.moreinfo}>
        {moreInfo ? (
          <MoreInfo
            setMoreInfo={setMoreInfo}
            data={data}
            certified={certified}
          />
        ) : null}
      </div>
    </div>
  );
}

export default RankDiv;
