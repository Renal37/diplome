import { useEffect, useState } from "react";
import "./admin_check_profile.css";

const AdminCheckProfile = () => {
    const [users, setUsers] = useState([]);
    const [searchQuery, setSearchQuery] = useState("");
    const [showFullInfo, setShowFullInfo] = useState(false);

    useEffect(() => {
        fetch("http://localhost:5000/users")
            .then((res) => res.json())
            .then((data) => setUsers(data));
    }, []);

    const filteredUsers = users.filter((user) =>
        user.username?.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div className="admin-delete-component">
            <h1>Проверка профилей</h1>
            <div className="controls">
                <div className="search-container">
                    <input
                        type="text"
                        placeholder="Поиск по имени пользователя"
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="search-input"
                    />
                    <button className="toggle-info-button" onClick={() => setShowFullInfo(!showFullInfo)}>
                        {showFullInfo ? "Скрыть детали" : "Показать детали"}
                    </button>
                </div>
            </div>
            <table className="course-table">
                <thead>
                    <tr>
                        <th>Логин</th>
                        <th>Email</th>
                        <th>ФИО</th>
                        {showFullInfo && (
                            <>
                                <th>Дата рождения</th>
                                <th>Место рождения</th>
                                <th>Образование</th>
                                <th>Место работы</th>
                                <th>Должность</th>
                                <th>Домашний адрес</th>
                                <th>Номер</th>
                                <th>Данные паспорта</th>
                                <th>СНИЛС</th>
                            </>
                        )}
                    </tr>
                </thead>
                <tbody>
                    {filteredUsers.length > 0 ? (
                        filteredUsers.map((user) => (
                            <tr key={user._id}>
                                <td>{user.username || "Данные отсутствуют"}</td>
                                <td>{user.email || "Данные отсутствуют"}</td>
                                <td>
                                    {user.lastname && user.firstname && user.middlename
                                        ? `${user.lastname} ${user.firstname} ${user.middlename}`
                                        : "Данные отсутствуют"}
                                </td>
                                {showFullInfo && (
                                    <>
                                        <td>{user.birthdate || "Данные отсутствуют"}</td>
                                        <td>{user.birthplace || "Данные отсутствуют"}</td>
                                        <td>{user.education || "Данные отсутствуют"}</td>
                                        <td>{user.workplace || "Данные отсутствуют"}</td>
                                        <td>{user.jobtitle || "Данные отсутствуют"}</td>
                                        <td>{user.homeaddress || "Данные отсутствуют"}</td>
                                        <td>{user.phone || "Данные отсутствуют"}</td>
                                        <td>{user.passportdata || "Данные отсутствуют"}</td>
                                        <td>{user.snils || "Данные отсутствуют"}</td>
                                    </>
                                )}
                            </tr>
                        ))
                    ) : (
                        <tr>
                            <td colSpan={showFullInfo ? 12 : 3} className="no-data-message">
                                Нет данных для отображения
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
    );
};

export default AdminCheckProfile;