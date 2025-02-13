import React, { useState, useEffect } from "react";
import "./admin_edit_group.css";

const AdminEditGroup = () => {
    const [groups, setGroups] = useState([]);
    const [selectedGroup, setSelectedGroup] = useState(null);
    const [groupName, setGroupName] = useState("");
    const [courseId, setCourseId] = useState("");
    const [error, setError] = useState("");
    const [successMessage, setSuccessMessage] = useState("");

    // Загрузка списка групп
    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch("http://localhost:5000/groups");
                if (!response.ok) {
                    throw new Error("Ошибка при загрузке групп");
                }
                const data = await response.json();
                setGroups(data.groups);
            } catch (error) {
                console.error("Error fetching groups:", error);
                setError("Ошибка при загрузке групп");
            }
        };
        fetchGroups();
    }, []);

    // Обработка редактирования группы
    const handleEdit = (group) => {
        setSelectedGroup(group);
        setGroupName(group.groupName);
        setCourseId(group.courseId);
    };

    // Обработка сохранения изменений
    const handleSave = async () => {
        if (!groupName || !courseId) {
            setError("Пожалуйста, заполните все поля");
            return;
        }

        try {
            const response = await fetch(
                `http://localhost:5000/admin/update-group/${selectedGroup._id}`,
                {
                    method: "PUT",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    credentials: "include",
                    body: JSON.stringify({ groupName, courseId }),
                }
            );

            if (!response.ok) {
                throw new Error("Ошибка при обновлении группы");
            }

            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Группа успешно обновлена!");
                setSelectedGroup(null);
                setGroupName("");
                setCourseId("");
                setError("");
                const updatedGroups = groups.map((g) =>
                    g._id === selectedGroup._id
                        ? { ...g, groupName, courseId }
                        : g
                );
                setGroups(updatedGroups);
            } else {
                setError(data.message || "Ошибка при обновлении группы");
            }
        } catch (error) {
            console.error("Error updating group:", error);
            setError("Ошибка при обновлении группы");
        }
    };

    // Обработка удаления группы
    const handleDelete = async (groupId) => {
        try {
            const response = await fetch(
                `http://localhost:5000/admin/delete-group/${groupId}`,
                {
                    method: "DELETE",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    credentials: "include",
                }
            );

            if (!response.ok) {
                throw new Error("Ошибка при удалении группы");
            }

            const data = await response.json();
            
            if (data.success) {
                setSuccessMessage("Группа успешно удалена!");
                setGroups(groups.filter((g) => g._id !== groupId));
            } else {
                setError(data.message || "Ошибка при удалении группы");
            }
        } catch (error) {
            console.error("Error deleting group:", error);
            setError("Ошибка при удалении группы");
        }
    };

    return (
        <div className="admin-groups-management-page">
            <h1>Управление группами</h1>
            {error && <div className="error-message">{error}</div>}
            {successMessage && (
                <div className="success-message">{successMessage}</div>
            )}
            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Название группы</th>
                        <th>ID курса</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {groups.map((group) => (
                        <tr key={group._id}>
                            <td>{group._id}</td>
                            <td>{group.groupName}</td>
                            <td>{group.courseId}</td>
                            <td>
                                <button onClick={() => handleEdit(group)}>Редактировать</button>
                                <button onClick={() => handleDelete(group._id)}>Удалить</button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>
            {selectedGroup && (
                <div>
                    <h2>Редактирование группы</h2>
                    <label>Название группы:</label>
                    <input
                        type="text"
                        value={groupName}
                        onChange={(e) => setGroupName(e.target.value)}
                    />
                    <label>ID курса:</label>
                    <input
                        type="text"
                        value={courseId}
                        onChange={(e) => setCourseId(e.target.value)}
                    />
                    <button onClick={handleSave}>Сохранить</button>
                </div>
            )}
        </div>
    );
};

export default AdminEditGroup;