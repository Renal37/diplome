import React, { useEffect, useState } from 'react';
import Header_center_head from "../main_center_head/main_center_head.component";
import "./main_top.component.css";

const Main_top = () => {
    const [headerHeight, setHeaderHeight] = useState(0);

    useEffect(() => {
        const headerBottom = document.querySelector('.header_bottom');
        if (headerBottom) {
            setHeaderHeight(headerBottom.offsetHeight);
        }
    }, []);

    return (
        <div className="main_top" style={{ paddingTop: headerHeight }}>
            <Header_center_head />
        </div>
    );
}

export default Main_top;