import React, { useState, useEffect } from "react";
import "./admin_group_management.css";

const AdminGroupManagement = () => {
    const [groups, setGroups] = useState([]);
    const [filteredGroups, setFilteredGroups] = useState([]); // Состояние для отфильтрованных групп
    const [searchTerm, setSearchTerm] = useState(""); // Состояние для строки поиска
    const [selectedGroup, setSelectedGroup] = useState(null);
    const [groupName, setGroupName] = useState("");
    const [courseId, setCourseId] = useState("");
    const [error, setError] = useState("");
    const [successMessage, setSuccessMessage] = useState("");
    const [isModalOpen, setIsModalOpen] = useState(false);

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
                setFilteredGroups(data.groups); // Инициализируем filteredGroups всеми группами
            } catch (error) {
                console.error("Error fetching groups:", error);
                setError("Ошибка при загрузке групп");
            }
        };
        fetchGroups();
    }, []);

    // Обработка изменения строки поиска
    const handleSearchChange = (e) => {
        const term = e.target.value.toLowerCase();
        setSearchTerm(term);

        // Фильтрация групп по строке поиска
        const filtered = groups.filter(
            (group) =>
                group.groupName.toLowerCase().includes(term) ||
                group.courseId.toString().toLowerCase().includes(term)
        );
        setFilteredGroups(filtered);
    };

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
                setFilteredGroups(updatedGroups); // Обновляем отфильтрованные группы
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
                setFilteredGroups(filteredGroups.filter((g) => g._id !== groupId)); // Обновляем отфильтрованные группы
            } else {
                setError(data.message || "Ошибка при удалении группы");
            }
        } catch (error) {
            console.error("Error deleting group:", error);
            setError("Ошибка при удалении группы");
        }
    };

    // Обработка создания группы
    const handleCreateGroup = async (e) => {
        e.preventDefault();
        if (!groupName || !courseId) {
            setError("Пожалуйста, заполните все поля");
            return;
        }
        try {
            const response = await fetch("http://localhost:5000/admin/create-group", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                credentials: "include",
                body: JSON.stringify({ groupName, courseId }),
            });
            const data = await response.json();
            if (data.success) {
                setSuccessMessage("Группа успешно создана!");
                setGroupName("");
                setCourseId("");
                setError("");
                setIsModalOpen(false);
                fetchGroups(); // Перезагружаем данные
            } else {
                setError(data.message || "Ошибка при создании группы");
            }
        } catch (error) {
            console.error("Error creating group:", error);
            setError("Ошибка при создании группы");
        }
    };

    // Закрытие модального окна и формы редактирования
    const handleClose = () => {
        setIsModalOpen(false);
        setSelectedGroup(null);
        setGroupName("");
        setCourseId("");
        setError("");
        setSuccessMessage("");
    };

    // Обработка нажатия Escape
    useEffect(() => {
        const handleKeyDown = (event) => {
            if (event.key === "Escape") {
                handleClose();
            }
        };
        window.addEventListener("keydown", handleKeyDown);
        return () => {
            window.removeEventListener("keydown", handleKeyDown);
        };
    }, []);

    return (
        <div className="admin-groups-management-page">
            {error && <div className="error-message">{error}</div>}
            {successMessage && <div className="success-message">{successMessage}</div>}

            <div className="admin-course-header">
                <input
                    type="text"
                    placeholder="Поиск по группам..."
                    value={searchTerm}
                    onChange={handleSearchChange}
                    className="search-input"
                />
                <button className="approve-btn" onClick={() => setIsModalOpen(true)}>
                    Создать группу
                </button>
            </div>

            {isModalOpen && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Создание новой группы</h2>
                        <form onSubmit={handleCreateGroup}>
                            <label>ID курса:</label>
                            <input
                                type="text"
                                value={courseId}
                                onChange={(e) => setCourseId(e.target.value)}
                                placeholder="Введите ID курса"
                            />
                            <label>Название группы:</label>
                            <input
                                type="text"
                                value={groupName}
                                onChange={(e) => setGroupName(e.target.value)}
                                placeholder="Введите название группы"
                            />
                            <div className="form-buttons">
                                <button className="approve-btn" type="submit">
                                    Создать
                                </button>
                                <button className="reject-btn" type="button" onClick={handleClose}>
                                    Закрыть
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}

            <table>
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Название группы</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    {filteredGroups.map((group) => (
                        <tr key={group._id}>
                            <td>{group.courseId}</td>
                            <td>{group.groupName}</td>
                            <td>
                                <button className="approve-btn" onClick={() => handleEdit(group)}>
                                    Редактировать
                                </button>
                                <button className="reject-btn" onClick={() => handleDelete(group._id)}>
                                    Удалить
                                </button>
                            </td>
                        </tr>
                    ))}
                </tbody>
            </table>

            {selectedGroup && (
                <div className="modal">
                    <div className="modal-content">
                        <h2>Редактирование группы</h2>
                        <form>
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
                            <div className="form-buttons">
                                <button type="button" onClick={handleSave}>
                                    Сохранить
                                </button>
                                <button type="button" onClick={handleClose}>
                                    Закрыть
                                </button>
                            </div>
                        </form>
                    </div>
                </div>
            )}
        </div>
    );
};

export default AdminGroupManagement;