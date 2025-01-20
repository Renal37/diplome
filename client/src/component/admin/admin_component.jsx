import AdminCourseManagement from "../admin_components/admin_add/admin_course_management";
import AdminUpdateComponent from "../admin_components/admin_update/admin_update_component";
import './admin_component.css';

const Admin = () => {
    return (
        <div className="admin">
            <AdminCourseManagement />
            <AdminUpdateComponent /> 
        </div>
    )
}

export default Admin;