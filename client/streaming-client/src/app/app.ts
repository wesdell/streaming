import { Component } from '@angular/core';

import { Home } from "./components/home/home";

@Component({
  selector: 'app-root',
  standalone: true,
  templateUrl: './app.html',
  imports: [Home]
})
export class App {

  handleReview(imdb: string) {
    console.log("Reviewing movie:", imdb);
  }

}
