import { Component, Input, signal, OnInit, inject } from '@angular/core';

import { Movies } from '../movies/movies';
import { Loader } from '../loader/loader';
import { MoviesApi } from '../../api/movie';

@Component({
  selector: 'home',
  standalone: true,
  imports: [Movies, Loader],
  templateUrl: './home.html'
})

export class Home implements OnInit {

  @Input() updateMovieReview!: (imdb: string) => void;

  movies = signal([]);
  loading = signal(false);
  message = signal('');

  private moviesApi = inject(MoviesApi);

  ngOnInit() {
    this.fetchMovies();
  }

  fetchMovies() {
    this.loading.set(true);
    this.message.set('');

    this.moviesApi.getMovies().subscribe({
      next: (data: any) => {
        this.movies.set(data);

        if (data.length === 0) {
          this.message.set('There are currently no movies available');
        }
      },
      error: () => {
        this.message.set('Error fetching movies');
      },
      complete: () => {
        this.loading.set(false);
      }
    });
  }
}
