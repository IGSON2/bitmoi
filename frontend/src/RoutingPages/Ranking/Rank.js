import { useEffect, useState } from "react";
import styles from "./Rank.module.css";
import RankDiv from "./RankDiv/RankDiv";
import Topbutton from "../../component/Topbutton/Topbutton";
import H_NavBar from "../../component/navbar/H_NavBar";
import axiosClient from "../../component/backendConn/axiosClient";
import AdDiv from "./AdDiv/AdDiv";
import mockup from "../../component/images/mockup_rank.png";
import getSelectedBidderImg from "../../component/backendConn/getSelectedBidderImg";
import Countdown from "./Countdown/Countdown";

function Rank() {
  const [pageNum, setPageNum] = useState(1);

  const [data, setData] = useState([{}]);
  const [imgLink, setImgLink] = useState("");
  const [nextUnlock, setNextUnlock] = useState();

  const getBidder = async () => {
    const adImgLink = await getSelectedBidderImg("rank");
    adImgLink ? setImgLink(adImgLink) : setImgLink(mockup);
  };

  const getUserScore = async () => {
    const response = await axiosClient.get(`/rank/${pageNum}`);
    setData(response.data);
  };

  const getNextBidUnlock = async () => {
    const res = await axiosClient.get("/nextBidUnlock");
    setNextUnlock(res.data.next_unlock);
  };

  useEffect(() => {
    getUserScore();
  }, [pageNum]);

  useEffect(() => {
    getNextBidUnlock();
    getBidder();
  }, []);

  return (
    <div className={styles.scorediv}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      <div className={styles.timer}>
        {nextUnlock ? <Countdown nextUnlock={nextUnlock} /> : null}
      </div>
      <div className={styles.graphbody}>
        <div className={styles.titlediv}>
          <h1 className={styles.title}>RANKING BOARD</h1>
        </div>
        {data
          ? data.map((v, i) => {
              if ((i + 1) % 10 === 0) {
                return (
                  <div className={styles.adline}>
                    <RankDiv key={i} index={i + 1} obj={v} />
                    <AdDiv imgLink={imgLink} />
                  </div>
                );
              }
              return <RankDiv key={i} index={i + 1} obj={v} />;
            })
          : null}

        <div className={styles.footer}>
          <div>Copyright &copy; 2023 IGSON All rights reserved.</div>
        </div>
      </div>
      <Topbutton />
    </div>
  );
}

export default Rank;
