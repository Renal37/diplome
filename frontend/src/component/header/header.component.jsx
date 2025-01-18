
import { Outlet } from 'react-router-dom';
import Header_top from '../header_components/header_top/header_top.component';
import Header_center from '../header_components/header_center/header_center.component';
import Header_bottom from '../header_components/header_bottom/header_bottom.component';
import { Fragment } from 'react';
import './header.component.css';

const Header = () => {
    return (
        <Fragment>
            <header>
                <div className="header container">
                    <Header_top />
                    <Header_center />
                </div>
                <Header_bottom />
            </header>
            <Outlet />
        </Fragment>
    )
}

export default Header;