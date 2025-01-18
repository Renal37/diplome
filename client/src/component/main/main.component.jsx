import { Outlet } from 'react-router-dom';
import Main_top from "../main_components/main_top/main_top.component";
import './main.component.css'
import Main_center from '../main_components/main_center/main.center.component';

const Main = () => {
    return (
        <div>
            <main className="container">
                <Main_top />
                <Outlet />
                <Main_center />
            </main>
        </div>
    );
};
export default Main;