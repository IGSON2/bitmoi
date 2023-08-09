import { useState } from "react";
import { Link } from "react-router-dom";
import styles from "./Home.module.css";
import HorizontalLine from "../../component/lines/HorizontalLine";
import HomeInformation from "./HomeInformation/HomeInformation";

function Home() {
  const menu = ["practice", "competition", "rank", "ad_bidding"];
  const [modetitle, setModetitle] = useState("BITMOI");
  const onMouseEnter = (event) => {
    const selected = event.target.outerText.toLowerCase();
    setModetitle(selected);
  };
  const onMouseLeave = () => {
    setModetitle("BITMOI");
  };
  return (
    <div className={styles.homepage}>
      <div className={styles.title}>
        <h6>바로 시작하는 모의투자</h6>
        <h1>BITMOI</h1>
      </div>
      <div className={`${styles.tophorizon} ${styles.hrs}`}>
        <HorizontalLine />
      </div>
      <div className={styles.informbox}>
        <div className={styles.informboxtitle}>
          <div>{modetitle.toUpperCase()}</div>
        </div>
        <HomeInformation mode={modetitle} />
      </div>
      <div className={`${styles.bottomhorizon}  ${styles.hrs}`}>
        <HorizontalLine />
      </div>
      <div className={styles.links}>
        {menu.map((m, idx) => {
          return (
            <Link to={`/${m}`} key={idx}>
              <button
                className={styles.linkbutton}
                onMouseEnter={onMouseEnter}
                onMouseLeave={onMouseLeave}
              >
                {m.toUpperCase().replace(/_/, " ")}
              </button>
            </Link>
          );
        })}
      </div>
      <div className={styles.copyright}>
        Copyright &copy; 2023 IGSON All rights reserved.
      </div>
    </div>
  );
}

export default Home;
