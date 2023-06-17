import styles from "./RankDiv.module.css";
import { BsChatQuote } from "react-icons/bs";
import { useEffect, useState } from "react";
import MoreInfo from "./MoreInfo/MoreInfo";
import { getAuth, onAuthStateChanged } from "firebase/auth";
import { BsForward } from "react-icons/bs";

function RankDiv({ index, obj }) {
  const [moreInfo, setMoreInfo] = useState(false);
  const [thisUser, setThisUser] = useState("");
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
      fetch(
        "http://localhost:5000/moreinfo/?user=" +
          obj.user +
          "&index=0&scoreid=" +
          obj.scoreid
      )
        .then((data) => {
          const json = data.json();
          return json;
        })
        .then((json) => {
          setData(json);
        });
    }

    setMoreInfo((current) => !current);
  };

  const auth = getAuth();
  // TODO: update userid for firebase
  useEffect(() => {
    onAuthStateChanged(auth, (user) => {
      setThisUser(user.uid);
    });
  }, []);
  useEffect(() => {
    if (thisUser && obj) {
      obj.user === thisUser ? setCertified(true) : setCertified(false);
    }
  }, [thisUser, obj]);
  return (
    <div
      className={styles.userdiv}
      style={obj.user === thisUser ? { backgroundColor: "#faebef" } : null}
    >
      <div className={styles.onlyinfo}>
        <div
          className={`${styles.no} ${styles.field}`}
          style={{ color: `${colorchanger(index)}` }}
        >
          {index}
        </div>
        <div className={`${styles.pic} ${styles.field}`}>
          <img className={styles.photo} src={obj.photourl} />
        </div>
        <div
          className={`${styles.name} ${styles.field}`}
          style={
            obj.user === thisUser
              ? {
                  color: "#333d79",
                  fontSize: "x-large",
                }
              : { color: "black" }
          }
        >
          {obj.displayname}
        </div>
        <div className={`${styles.score}  ${styles.field}`}>{obj.balance}</div>
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
