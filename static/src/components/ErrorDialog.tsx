import Dialog from 'material-ui/Dialog';
import FlatButton from 'material-ui/FlatButton';
import * as React from 'react';

export interface IErrorDialogProps {
  open: boolean;
  error: string;
  onRequestClose: () => void;
}

export class ErrorDialog extends React.Component<IErrorDialogProps, any> {
  public render() {
    const actions = [<FlatButton key="OK" label="OK" primary={true} onClick={this.props.onRequestClose} />]; // FIXME

    return (
      <Dialog
        title="Error occurred"
        actions={actions}
        modal={false}
        open={this.props.open}
        onRequestClose={this.props.onRequestClose}
      >
        {this.props.error}
      </Dialog>
    );
  }
}
