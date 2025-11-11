import { Component } from '@angular/core';
import { Movie } from './components/movie/movie';

@Component({
  selector: 'app-root',
  standalone: true,
  templateUrl: './app.html',
  imports: [Movie]
})

export class App {
  movie = {
    _id: '1',
    title: 'Interestelar',
    imdb_id: 'tt0816692',
    youtube_id: 'zSWdZVtXT7E',
    poster_path: 'https://image.tmdb.org/t/p/w500/nBNZadXqJSdt05SHLqgT0HuC5Gm.jpg',
    ranking: {
      ranking_name: 'Top Rated'
    }
  };

  handleReview(imdb: string) {
    console.log('Reviewing movie:', imdb);
  }
}
