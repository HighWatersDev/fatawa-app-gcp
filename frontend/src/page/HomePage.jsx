import React, {useEffect, useState} from 'react';
import { Link } from 'react-router-dom';
import { onAuthStateChanged } from "firebase/auth";
import { auth } from '../firebase';
import styles from '../components/styles/HomePage.module.css'
import { FaCopy } from "react-icons/fa";
import { CopyToClipboard } from "react-copy-to-clipboard";

const HomePage = () => {
    const [userToken, setUserToken] = useState(null);
    const [isCopied, setIsCopied] = useState(false);

    useEffect(()=>{
        onAuthStateChanged(auth, (user) => {
            if (user) {
                user.getIdToken().then(token => setUserToken(token));

            } else {

                console.log("user is logged out")
            }
        })})
    const handleCopy = () => {
        setIsCopied(true);
        setTimeout(() => {
            setIsCopied(false);
        }, 3000);
    };

    return (
        <div>
            <h1>Welcome to Salafi Fatawa App</h1>
            <nav>
                <ul>
                    <li>
                        <Link to="/login">Login</Link>
                    </li>
                    <li>
                        <Link to="/fatawa">Fatawa</Link>
                    </li>
                </ul>
            </nav>
            <div className={styles.container}>
                <div className={styles.tokenBox}>
                    {userToken && (
                        <code className={styles.tokenCode}>
                            <span className={styles.tokenLabel}>Your token:</span>{" "}
                            {userToken}
                        </code>
                    )}
                </div>
                <CopyToClipboard text={userToken} onCopy={handleCopy}>
                    <button className={styles.copyButton}>
                        <FaCopy className={styles.copyIcon} />
                        {isCopied ? "Copied!" : "Copy token"}
                    </button>
                </CopyToClipboard>
            </div>
        </div>
    );
}

export default HomePage;