import React from 'react';
import Slider from 'react-slick';
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import './main_top_center.component.css';
import apt from "../../../../../assets/apt.png";

const Main_top_center = () => {
    const settings = {
        dots: true,
        infinite: true,
        speed: 500,
        slidesToShow: 1,
        slidesToScroll: 1,
        autoplay: true,
        autoplaySpeed: 3000,
    };

    return (
        <div className="main_top_center_slider">
           <img src={apt} alt="" />
                <Slider {...settings} className='main_top_slider'>
                    <div>
                        <h3>Почему наш сайт?</h3>
                        <p>Мы предлагаем качественные образовательные программы, разработанные опытными преподавателями.</p>
                    </div>
                    <div>
                        <h3>Почему именно у нас?</h3>
                        <p>Наши курсы помогут вам получить актуальные знания и навыки, необходимые для успешной карьеры.</p>
                    </div>
                    <div>
                        <h3>Преимущества наших курсов</h3>
                        <p>Гибкий график обучения, доступные цены и поддержка на всех этапах обучения.</p>
                    </div>
                </Slider>
              
            </div>
    );
}

export default Main_top_center;