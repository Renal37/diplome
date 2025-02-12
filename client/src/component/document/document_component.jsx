import React from "react";
import docmuent1 from "../../assets/med.jpg";
import docmuent2 from "../../assets/lic.jpg";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faFilePdf } from "@fortawesome/free-solid-svg-icons";
import "./document_component.css";

const Document = () => {
    return (
        <div className="document-container">
            <div className="image-container">
                <img src={docmuent1} alt="Устав" className="document-image" />
                <img src={docmuent2} alt="Лицензия" className="document-image" />
            </div>
            <div className="text-container">
                <div className="document_card">
                    <p className="document-info">
                        УСТАВ Государственного Автономного Профессионального Образовательного Учреждения «Альметьевский Политехнический Техникум» от 10.11.22г. №1911/22
                    </p>
                    <a href="https://almetpt.ru/files/about/npb/doc/ustav.pdf" className="icon-link">
                        <FontAwesomeIcon icon={faFilePdf} style={{ color: "#ffffff", }} />
                    </a>
                </div>
                <div className="document_card">
                    <p className="document-info">
                        Изменения в УСТАВ Государственного Автономного Профессионального Образовательного Учреждения «Альметьевский Политехнический Техникум» от 17.11.23г. №2064/23
                    </p>
                    <a href="https://almetpt.ru/files/about/npb/doc/ustav_izm.pdf" className="icon-link">
                        <FontAwesomeIcon icon={faFilePdf} style={{ color: "#ffffff", }} />
                    </a>
                </div>
                <div className="document_card">
                    <p className="document-info">
                        Лицензия на осуществление образовательной деятельности от 03 июля 2015 г. № 6727
                    </p>
                    <a href="https://almetpt.ru/img/html/about/npb/doc/licenziya.pdf" className="icon-link">
                        <FontAwesomeIcon icon={faFilePdf} style={{ color: "#ffffff", }} />
                    </a>
                </div>
                <div className="document_card">
                    <p className="document-info">
                        Выписка из реестра лицензий № Л035-01272-16/00252330 от 03 июля 2015г. (Лицензия от 03 июля 2015 г. № 6727)
                    </p>
                    <a href="https://almetpt.ru/files/about/npb/doc/vypiska.pdf" className="icon-link">
                        <FontAwesomeIcon icon={faFilePdf} style={{ color: "#ffffff", }} />
                    </a>
                </div>
                <div className="document_card">
                    <p className="document-info">
                        Выписка из государственной информационной системы «Реестр организаций, осуществляющих образовательную деятельность по имеющим государственную аккредитацию образовательным программам» №3190 от 11.11.2015
                    </p>
                    <a href="https://almetpt.ru/files/about/npb/doc/akkred.pdf" className="icon-link">
                        <FontAwesomeIcon icon={faFilePdf} style={{ color: "#ffffff", }} />
                    </a>
                </div>
                <div className="document_card">
                    <p className="document-info">
                        Регистрация в реестре организаций, оказывающих услуги в области охраны труда (регистрационный номер №9491 от 02.04.2024г.)
                    </p>
                    <a href="https://almetpt.ru/files/about/npb/doc/reestr_ot.pdf" className="icon-link">
                        <FontAwesomeIcon icon={faFilePdf} style={{ color: "#ffffff", }} />
                    </a>
                </div>
            </div>
        </div>
    );
};

export default Document;    