import { useEffect, useState } from "react";

const AdminCheckProfile = () => {
    const [users, setUsers] = useState([]);
    const [searchQuery, setSearchQuery] = useState("");

    useEffect(() => {
        fetch("http://localhost:5000/users")
            .then(res => res.json())
            .then(data => setUsers(data));
    }, []);

    const filteredUsers = users.filter(user =>
        user.username.toLowerCase().includes(searchQuery.toLowerCase())
    );

    return (
        <div>
            <h1>Проверка профилей</h1>
            <input
                type="text"
                placeholder="Поиск по имени пользователя"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full p-2 mb-4 border rounded"
            />
            {filteredUsers.map(user => (
                <div key={user._id}>
                    <p>ЛОГИН:   {user.username}</p>
                    <p>ФИО: {user.fullName}</p>
                </div>
            ))}
        </div>
    );
};

export default AdminCheckProfile;