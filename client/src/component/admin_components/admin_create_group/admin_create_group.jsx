import React, { useState } from "react";
import "./admin_create_group.css";

const AdminCreateGroup = () => {
    const [groupName, setGroupName] = useState("");
    const [courseId, setCourseId] = useState("");
    const [error, setError] = useState("");
    const [successMessage, setSuccessMessage] = useState("");

    const handleSubmit = async (e) => {
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
            } else {
                setError(data.message || "Ошибка при создании группы");
            }
        } catch (error) {
            console.error("Error creating group:", error);
            setError("Ошибка при создании группы");
        }
    };

    return (
        <div className="admin-group-creation-page">
            <h1>Создание новой группы</h1>
            {error && <div className="error-message">{error}</div>}
            {successMessage && <div className="success-message">{successMessage}</div>}
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Название группы:</label>
                    <input
                        type="text"
                        value={groupName}
                        onChange={(e) => setGroupName(e.target.value)}
                        placeholder="Введите название группы"
                    />
                </div>
                <div>
                    <label>ID курса:</label>
                    <input
                        type="text"
                        value={courseId}
                        onChange={(e) => setCourseId(e.target.value)}
                        placeholder="Введите ID курса"
                    />
                </div>
                <button type="submit">Создать группу</button>
            </form>
        </div>
    );
};

export default AdminCreateGroup;