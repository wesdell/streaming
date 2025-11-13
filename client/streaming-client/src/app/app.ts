import { Component } from '@angular/core';

import { Home } from "./components/home/home";
import { Header } from "./components/header/header";

@Component({
  selector: 'app-root',
  standalone: true,
  templateUrl: './app.html',
  imports: [Home, Header]
})
export class App {

  handleReview(imdb: string) {
    console.log("Reviewing movie:", imdb);
  }

}
