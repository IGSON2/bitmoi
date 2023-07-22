import { useParams } from "react-router-dom";
import styles from "./EmailLocate.module.css";
import H_NavBar from "../../component/navbar/H_NavBar";

function EmailLocate() {
  const { domain } = useParams();
  const emailPlatform = domain.split(".")[0];
  const targetURL = `https://${domain}`;

  const UpperDomain = `${emailPlatform
    .slice(0, 1)
    .toUpperCase()}${emailPlatform.slice(1, emailPlatform.length)}`;
  return (
    <div className={styles.page}>
      <div className={styles.navbar}>
        <H_NavBar />
      </div>
      <h1 className={styles.message}>이메일 인증을 완료해 주세요</h1>
      <a className={styles.link} href={targetURL} target="_blank">
        <button>{UpperDomain}로 이동하기</button>
      </a>
    </div>
  );
}

export default EmailLocate;
