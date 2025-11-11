import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';

import { Movie } from '../movie/movie';
import { IMovie } from '../../interfaces/movie';

@Component({
  selector: 'movies',
  standalone: true,
  imports: [CommonModule, Movie],
  templateUrl: './movies.html',
})
export class Movies {
  @Input() movies: IMovie[] = [];
  @Input() message: string = '';

  @Output() updateMovieReview = new EventEmitter<string>();

  emitUpdate(movieId: string) {
    this.updateMovieReview.emit(movieId);
  }
}
