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
  @Output() updateMovieReview = new EventEmitter<string>();

  emitReview() {
    this.updateMovieReview.emit(this.movie.imdb_id);
  }
}
