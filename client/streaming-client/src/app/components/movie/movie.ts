import { Component, Input, Output, EventEmitter } from '@angular/core';
import { RouterModule } from '@angular/router';
import { CommonModule } from '@angular/common';

import { IMovie } from '../../interfaces/movie';

@Component({
  selector: 'movie',
  standalone: true,
  templateUrl: './movie.html',
  imports: [RouterModule, CommonModule]
})

export class Movie {
  @Input() movie!: IMovie;
  @Input() updateMovieReview!: (id: string) => void;

  emitReview() {
    this.updateMovieReview(this.movie.imdb_id);
  }
}
