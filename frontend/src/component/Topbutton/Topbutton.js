import { BsArrowUpCircle } from "react-icons/bs";
import styles from "./Topbutton.module.css";
function Topbutton() {
  const goTop = () => {
    document.documentElement.scrollTop = 0;
  };

  return (
    <div className={styles.topbutton}>
      <button onClick={goTop}>
        <BsArrowUpCircle />
      </button>
    </div>
  );
}

export default Topbutton;
