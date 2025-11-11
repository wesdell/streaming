import { Component } from '@angular/core';

import { Movies } from './components/movies/movies';

@Component({
  selector: 'app-root',
  standalone: true,
  templateUrl: './app.html',
  imports: [Movies]
})

export class App {
  movies = [
    {
      _id: '1',
      title: 'Interestelar',
      imdb_id: 'tt0816692',
      youtube_id: 'zSWdZVtXT7E',
      poster_path: 'https://image.tmdb.org/t/p/w500/nBNZadXqJSdt05SHLqgT0HuC5Gm.jpg',
      ranking: {
        ranking_name: 'Top Rated'
      }
    },
    {
      _id: '2',
      title: 'Interestelar',
      imdb_id: 'tt0816692',
      youtube_id: 'zSWdZVtXT7E',
      poster_path: 'https://image.tmdb.org/t/p/w500/nBNZadXqJSdt05SHLqgT0HuC5Gm.jpg',
      ranking: {
        ranking_name: 'Top Rated'
      }
    }
  ];

  message = "No movies to show";

  handleReview(imdb: string) {
    console.log('Reviewing movie:', imdb);
  }
}
